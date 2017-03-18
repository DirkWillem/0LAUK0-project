import { Model, SubModelClass, PartialModel, ModelJson } from "../model";
import { AuthHttp } from "./authhttp.service";

/**
 * Provides the base service for all API interface services
 */
export class APIInterface<M extends Model> {
  protected baseURL: string;
  protected model: SubModelClass<M>;

  constructor(protected http: AuthHttp) {

  }

  /**
   * Creates a new entity
   * @param data - The data of the entity to be created
   * @returns {Promise<M>} - The created entity
   */
  async create(data: any): Promise<M> {
    const json = await this.http.postJSON<ModelJson<M>>(`/api${this.baseURL}`, data);
    return new this.model(json);
  }

  /**
   * Returns a list of all entities of the API from the server
   * @returns {Promise<M[]>} promise resolving to the list of entities
   */
  async list(): Promise<M[]> {
    const json = await this.http.getJSON<ModelJson<M>[]>(`/api${this.baseURL}`);
    return json.map(item => new this.model(item));
  }

  /**
   * Returns a list single entity by its ID
   * @returns {Promise<M[]>} promise resolving to the read entity
   */
  async read(id: number): Promise<M> {
    const json = await this.http.getJSON<ModelJson<M>>(`/api${this.baseURL}/${id}`);
    return new this.model(json);
  }

  /**
   * Deletes an entity
   * @param id - The ID of the entity to delete
   */
  async delete(id: number): Promise<void> {
    await this.http.delete(`/api${this.baseURL}/${id}`);
  }
}