import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';
import * as equal from 'deep-equal';

/**
 * Contains the data in a message returned by the dispatcher
 */
export interface DispatcherMessage<T> {
  subscriptionId: number;
  requestId: number;
  action: string;
  payload: T;
}

/**
 * Contains a request that is sent to the dispatcher
 */
interface DispatcherRequest<T> {
  requestId: number;
  action: string;
  payload: T;
}

/**
 * Contains the payload for a subscription request
 */
interface SubscribeRequestPayload {
  subject: string;
  subscriptionParams: any;
}

/**
 * Contains the payload for an unsubscription request
 */
interface UnsubscribeRequestPayload {
  subscriptionId: number;
}

/**
 * Contains the payload returned by a subscribe response
 */
interface SubscribeMessagePayload {
  subscriptionId: number;
}

/**
 * Contains information on a dispatcher subscription
 */
export interface DispatcherSubscription<T> {
  updates: Observable<T>;
  subscriptionId: number;
}

/**
 * Internal reference of a subject subscription in the dispatcher service
 */
interface SubjectSubscription {
  subscription: DispatcherSubscription<any>;
  subject: string;
  subscriptionParams: any;
  subscriberCount: number;
}

/**
 * Service for interfacing with the dispatcher
 */
@Injectable()
export class DispatcherService {
  socket: WebSocket;

  private messages: Subject<DispatcherMessage<any>> = new Subject<DispatcherMessage<any>>();
  private requestIdCounter: number = 1;
  private connection: Promise<void>;

  private subscriptions: SubjectSubscription[] = [];

  constructor() {
    this.connection = new Promise<void>(resolve => {
      this.socket = new WebSocket(`ws://${window.location.host == 'localhost' ? 'localhost:5000' : window.location.host}/api/dispatcher`);
      this.socket.addEventListener("open", (e) => {
        resolve();
      });

      this.socket.addEventListener("message", (e) => {
        this.messages.next(JSON.parse(e.data));
      });
    });

  }

  /**
   * Returns a subscription to a given subject with the given parameters
   * @param subject - The subject to subscribe to
   * @param subscriptionParams - The parameters to subscribe with
   * @returns {Promise<DispatcherSubscription<T>>} Promise resolving to the subscription observable
   */
  async subscribeTo<T>(subject: string, subscriptionParams: any): Promise<DispatcherSubscription<DispatcherMessage<T>>> {
    // Check if there is an existing subscription already
    const existing = this.subscriptions.find(sub => sub.subject == subject && equal(sub.subscriptionParams, subscriptionParams));

    if(existing) {
      this.subscriptions.forEach(sub => {
        if(sub.subscription.subscriptionId == existing.subscription.subscriptionId) {
          sub.subscriberCount++;
        }
      });

      return existing.subscription;
    }

    // Otherwise, create a new subscription
    const response = await this.sendRequest<SubscribeRequestPayload, SubscribeMessagePayload>("subscribe", {subject, subscriptionParams});

    console.log(response);

    const subscription = {
      updates: this.messages
        .filter(msg => msg.subscriptionId == response.payload.subscriptionId),
      subscriptionId: response.payload.subscriptionId
    };

    this.subscriptions.push({subscription, subject, subscriptionParams, subscriberCount: 1});
    return subscription;
  }

  async unsubscribeTo(subscriptionId: number) {
    await Promise.all(this.subscriptions.map(async (sub) => {
      if(sub.subscription.subscriptionId == subscriptionId) {
        sub.subscriberCount--;
        if(sub.subscriberCount <= 0) {
          await this.sendRequest<UnsubscribeRequestPayload, {}>("unsubscribe", {subscriptionId})
        }
      }

      return sub;
    }));

    this.subscriptions = this.subscriptions.filter(sub => sub.subscriberCount > 0);
  }

  /**
   * Sends a request to the
   * @param action
   * @param payload
   * @returns {Promise<number>}
   */
  private async sendRequest<RqP, RsP>(action: string, payload: RqP): Promise<DispatcherMessage<RsP>> {
    await this.connection;
    this.socket.send(JSON.stringify({
      requestId: this.requestIdCounter,
      action,
      payload
    } as DispatcherRequest<RqP>));

    const requestId = this.requestIdCounter++;

    return await new Promise<DispatcherMessage<RsP>>(resolve => {
      const subscription = this.messages
        .filter(msg => msg.subscriptionId == 0 && msg.requestId == requestId)
        .subscribe(msg => {
          subscription.unsubscribe();
          resolve(msg as DispatcherMessage<RsP>);
        });
    });
  }
}