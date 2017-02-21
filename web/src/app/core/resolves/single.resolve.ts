import { Resolve, ActivatedRouteSnapshot } from "@angular/router";

import { Model } from "../model";
import { APIInterface } from "../services/apiinterface.service";

/**
 * Base resolve for a single entity through an APIInterface
 */
export class SingleResolve<M extends Model> implements Resolve<M> {
  constructor(private apiService: APIInterface<M>) { }

  async resolve(route: ActivatedRouteSnapshot): Promise<M> {
    let r = route;
    while(!r.params["id"]) {
      r = r.parent;
    }

    const id = r.params["id"];

    return await this.apiService.read(parseInt(id));
  }
}