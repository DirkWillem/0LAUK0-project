import { Injectable } from "@angular/core";

import { NestedListResolve } from "./list.resolve";
import { PRNMedicationService, PRNMedication } from "../services/prnmedication.service";

/**
 * Resolver for a list of PRN medications of a user
 */
@Injectable()
export class PRNMedicationsResolve extends NestedListResolve<PRNMedication> {
  constructor(service: PRNMedicationService) {
    super(service);
  }
}