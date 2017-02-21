import { Model, Field, ModelListField } from "../model";
import { Injectable } from "@angular/core";
import { APIInterface } from "./apiinterface.service";
import { AuthHttp } from "./authhttp.service";

/**
 * Model for a user
 */
export class User extends Model {
  @Field() id: number;
  @Field() username: string;
  @Field() fullName: string;
  @Field() role: string;
  @Field() email: string;
  @Field() emailMD5: string;

  @ModelListField({optional: true, detail: true, model: User}) doctors: User[];
  @ModelListField({optional: true, detail: true, model: User}) pharmacists: User[];
  @ModelListField({optional: true, detail: true, model: User}) patients: User[];
  @ModelListField({optional: true, detail: true, model: User}) customers: User[];
}

@Injectable()
export class UserService extends APIInterface<User> {
  baseURL = "/users";
  model = User;

  constructor(authHttp: AuthHttp) {
    super(authHttp);
  }
}