import { Model, SubModelClass, PartialModel, ModelJson } from "../model";
import { AuthHttp } from "./authhttp.service";

/**
 * Provides the base service for all nested API interface services
 */
export class NestedAPIInterface<M extends Model> {
  protected baseURL: string;
  protected nestedURL: string;
  protected model: SubModelClass<M>;

  constructor(protected http: AuthHttp) {

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
}