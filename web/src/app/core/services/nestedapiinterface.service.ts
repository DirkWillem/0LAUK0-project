import { Observable } from 'rxjs';

import { Model, SubModelClass, PartialModel, ModelJson } from "../model";
import { AuthHttp } from "./authhttp.service";
import { DispatcherService, DispatcherSubscription } from "./dispatcher.service";
import { CollectionMutation, CollectionAddition, CollectionUpdate, CollectionRemoval } from "../collectionupdates";

/**
 * Provides the base service for all nested API interface services
 */
export class NestedAPIInterface<M extends Model> {
  protected baseURL: string;
  protected nestedURL: string;
  protected model: SubModelClass<M>;

  protected collectionSubject: string = "";
  protected collectionSubjectSuperIdProperty: string = "superId";

  constructor(protected http: AuthHttp, protected dispatcherService: DispatcherService) {

  }

  /**
   * Creates a new entity
   * @param superId - ID of the super-entity
   * @param data - The data of the entity to be created
   * @returns {Promise<M>} - The created entity
   */
  async create(superId: number, data: any): Promise<M> {
    const json = await this.http.postJSON<ModelJson<M>>(`/api${this.baseURL}/${superId}${this.nestedURL}`, data);
    return new this.model(json);
  }

  /**
   * Returns a list of all entities of the API from the server
   * @param superId - ID of the super-entity
   * @returns {Promise<M[]>} promise resolving to the list of entities
   */
  async list(superId: number): Promise<M[]> {
    const json = await this.http.getJSON<ModelJson<M>[]>(`/api${this.baseURL}/${superId}${this.nestedURL}`);
    return json.map(item => new this.model(item));
  }

  /**
   * Returns a list single entity by its ID
   * @param superId - ID of the super-entity
   * @param id - ID of the entity to read
   * @returns {Promise<M[]>} promise resolving to the read entity
   */
  async read(superId: number, id: number): Promise<M> {
    const json = await this.http.getJSON<ModelJson<M>>(`/api${this.baseURL}/${superId}${this.nestedURL}/${id}`);
    return new this.model(json);
  }

  /**
   * Updates an entity by its ID
   * @param superId - ID of the super-entity
   * @param id - ID of the entity to read
   * @param model - Updated entity value
   * @returns {Promise<M>} Promise resolving to the updated entity
   */
  async update(superId: number, id: number, model: M) {
    const json = await this.http.putJSON<ModelJson<M>[]>(`/api${this.baseURL}/${superId}${this.nestedURL}/${id}`, model.toJSON());
    return new this.model(json);
  }

  /**
   * Deletes an entity by its ID
   * @param superId - The ID of the super entity of the entity to delete
   * @param id - The ID of the entity to delete
   * @returns {Promise<void>} Promise that resolves once the delete is done
   */
  async delete(superId: number, id: number): Promise<void> {
    return await this.http.delete(`/api${this.baseURL}/${superId}${this.nestedURL}/${id}`)
  }

  /**
   * Returns an observable which contains the collection updates returned by the dispatcher
   * @param superId - ID of the super entity to subscribe to
   */
  async getCollectionUpdates(superId: number): Promise<DispatcherSubscription<CollectionMutation<M>>> {
    const sub = (await this.dispatcherService
      .subscribeTo(this.collectionSubject, {
        [this.collectionSubjectSuperIdProperty]: superId
      }));

    return {
      updates: sub.updates
        .map(msg => {
        console.log(msg);
        switch(msg.action) {
          case "added":
            return <CollectionAddition<M>>{
              action: "added",
              addedEntity: new this.model((msg.payload as {addedEntity: ModelJson<M>}).addedEntity)
            };
          case "updated":
            return <CollectionUpdate<M>> {
              action: "updated",
              updatedEntity: new this.model((msg.payload as {updatedEntity: ModelJson<M>}).updatedEntity)
            };
          case "deleted":
            return <CollectionRemoval> {
              action: "deleted",
              id: (msg.payload as {id: number}).id
            };
        }
      }),
      subscriptionId: sub.subscriptionId
    }
  }
}