import { Injectable } from '@angular/core';

import { NestedAPIInterface } from "./nestedapiinterface.service";
import { Model, Field, ModelListField } from "../model";
import { Medication } from "./medication.service";
import { AuthHttp } from "./authhttp.service";

/**
 * Class representing a dose
 */
export class Dose extends Model {
  @Field() id: number;
  @Field() title: string;
  @Field() dispenseBefore: string;
  @Field() dispenseAfter: string;
  @Field({detail: true}) description: string;
  @ModelListField({detail: true, model: Medication}) medications: Medication[];
}

export interface NewDose {
  title: string;
  dispenseAfter: {hour: number, minute: number};
  dispenseBefore: {hour: number, minute: number};
  description: string;
  medications: {medicationId: number, amount: number}[];
}

/**
 * Service for interfacing with the doses API
 */
@Injectable()
export class DoseService extends NestedAPIInterface<Dose> {
  baseURL = "/users";
  nestedURL = "/doses";

  model = Dose;

  constructor(authHttp: AuthHttp) {
    super(authHttp);
  }

  async create(superId: number, newDose: NewDose): Promise<Dose> {
    let dose: any = newDose;
    dose.dispenseAfter = `${newDose.dispenseAfter.hour.toString()}:${newDose.dispenseBefore.minute.toString()}:00`;
    dose.dispenseBefore = `${newDose.dispenseBefore.hour.toString()}:${newDose.dispenseBefore.minute.toString()}:00`;

    return await super.create(superId, dose);
  }
}