import { Injectable } from '@angular/core';
import { APIInterface } from "./apiinterface.service";
import { Model, Field } from "../model";
import { AuthHttp } from "./authhttp.service";

/**
 * Model representing a medication
 */
export class Medication extends Model {
  @Field() id: number;
  @Field() title: string;
  @Field() description: string;
}

/**
 * Service for interfacing with the medications API
 */
@Injectable()
export class MedicationService extends APIInterface<Medication> {
  baseURL = "/medications";
  model = Medication;

  constructor(authHttp: AuthHttp) {
    super(authHttp);
  }
}