import { Injectable } from '@angular/core';
import { APIInterface } from "./apiinterface.service";
import { Model, Field } from "../model";
import { AuthHttp } from "./authhttp.service";

export class Medication extends Model {
  @Field() id: number;
  @Field() title: string;
  @Field() description: string;
}

@Injectable()
export class MedicationService extends APIInterface<Medication> {
  baseURL = "/medications";
  model = Medication;

  constructor(authHttp: AuthHttp) {
    super(authHttp);
  }
}