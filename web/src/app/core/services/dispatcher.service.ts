import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';

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
 * Contains the payload returned by a subscribe response
 */
interface SubscribeMessagePayload {
  subscriptionId: number;
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

  constructor() {
    this.connection = new Promise<void>(resolve => {
      this.socket = new WebSocket(`ws://localhost:5000/api/dispatcher`);
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
   * @returns {Promise<Observable<DispatcherMessage<T>>>} Promise resolving to the subscription observable
   */
  async subscribeTo<T>(subject: string, subscriptionParams: any): Promise<Observable<DispatcherMessage<T>>> {
    const response = await this.sendRequest<SubscribeRequestPayload, SubscribeMessagePayload>("subscribe", {subject, subscriptionParams});

    return this.messages
      .filter(msg => msg.subscriptionId == response.payload.subscriptionId);
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