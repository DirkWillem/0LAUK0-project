import { Injectable } from "@angular/core";

import { ListResolve } from "./list.resolve";
import { Medication, MedicationService } from "../services/medication.service";
import { SingleResolve } from "./single.resolve";

/**
 * Resolver for a list of medications
 */
@Injectable()
export class MedicationsResolve extends ListResolve<Medication> {
  constructor(service: MedicationService) {
    super(service);
  }
}

/**
 * Resolver for a single medication
 */
@Injectable()
export class MedicationResolve extends SingleResolve<Medication> {
  constructor(service: MedicationService) {
    super(service);
  }
}