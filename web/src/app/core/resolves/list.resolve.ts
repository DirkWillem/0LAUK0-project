import { Resolve, ActivatedRouteSnapshot } from "@angular/router";

import { Model } from "../model";
import { APIInterface } from "../services/apiinterface.service";
import { NestedAPIInterface } from "../services/nestedapiinterface.service";

/**
 * Base resolve for a list of entities through an APIInterface
*/
export class ListResolve<M extends Model> implements Resolve<M[]> {
  constructor(protected apiService: APIInterface<M>) { }

  async resolve(): Promise<M[]> {
    return await this.apiService.list();
  }
}

/**
 * Base resolve for a list of nested entities through an NestedAPIInterface
 */
export class NestedListResolve<M extends Model> implements Resolve<M[]> {
  protected superIdProperty: string = "id";


  constructor(protected apiService: NestedAPIInterface<M>) { }

  async resolve(route: ActivatedRouteSnapshot): Promise<M[]> {
    let r = route;
    while(!r.params[this.superIdProperty]) {
      r = r.parent;
    }

    const id = r.params[this.superIdProperty];

    return await this.apiService.list(parseInt(id));
  }
}