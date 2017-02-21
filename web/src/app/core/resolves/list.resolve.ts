import { Resolve } from "@angular/router";

import { Model } from "../model";
import { APIInterface } from "../services/apiinterface.service";

/**
 * Base resolve for a list of entities through an APIInterface
 */
export class ListResolve<M extends Model> implements Resolve<M[]> {
  constructor(private apiService: APIInterface<M>) { }

  async resolve(): Promise<M[]> {
    return await this.apiService.list();
  }
}