import { Injectable } from '@angular/core';

import { NestedAPIInterface } from "./nestedapiinterface.service";
import { Model, Field, ModelField } from "../model";
import { Medication } from "./medication.service";
import { AuthHttp } from "./authhttp.service";
import { DispatcherService } from "./dispatcher.service";

/**
 * Model containing a PRN medication
 */
export class PRNMedication extends Model {
  @Field() description: string;
  @Field() maxDaily: number;
  @Field() minInterval: number;
  @Field() userId: number;
  @ModelField({model: Medication}) medication: Medication;
}

/**
 * Interface containing the data sent to the service on the creation of a PRN medication
 */
export interface NewPRNMedication {
  description: string;
  maxDaily: number;
  minInterval: number;
  medicationId: number;
}

/**
 * Service for interacting with the PRN medications API
 */
@Injectable()
export class PRNMedicationService extends NestedAPIInterface<PRNMedication> {
  baseURL = "/users";
  nestedURL = "/prnmedications";

  model = PRNMedication;

  protected collectionSubject: string = "doses";
  protected collectionSubjectSuperIdProperty: string = "userId";

  constructor(authHttp: AuthHttp, dispatcherService: DispatcherService) {
    super(authHttp, dispatcherService);
  }

  async create(superId: number, newPRNMedication: NewPRNMedication) {
    return await super.create(superId, newPRNMedication);
  }
}