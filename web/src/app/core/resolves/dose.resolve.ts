import { Injectable } from "@angular/core";

import { NestedListResolve } from "./list.resolve";
import { DoseService, Dose } from "../services/dose.service";

/**
 * Resolver for a list of doses of a user
 */
@Injectable()
export class DosesResolve extends NestedListResolve<Dose> {
  constructor(service: DoseService) {
    super(service);
  }
}