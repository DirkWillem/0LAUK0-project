import { Model, Field, ModelListField, DateField, ModelJson } from "../model";
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
  @Field() phone: string;
  @Field({detail: true}) gender: string;
  @DateField({detail: true}) birthdate: Date;

  @ModelListField({optional: true, detail: true, model: User}) doctors: User[];
  @ModelListField({optional: true, detail: true, model: User}) pharmacists: User[];
  @ModelListField({optional: true, detail: true, model: User}) patients: User[];
  @ModelListField({optional: true, detail: true, model: User}) customers: User[];
}

/**
 * Data sent to the service for creating a new user
 */
export interface NewUser {
  username: string;
  fullName: string;
  role: string;
  email: string;
  birthdate: Date;
  gender: string;
  phone: string;

  patientIds?: number[];
  customerIds?: number[];
  doctorIds?: number[];
  pharmacistIds?: number[];
}

@Injectable()
export class UserService extends APIInterface<User> {
  baseURL = "/users";
  model = User;

  constructor(authHttp: AuthHttp) {
    super(authHttp);
  }

  /**
   * Creates a new user
   * @param user - The new user data
   * @returns {Promise<M>} the created user
   */
  async create(user: NewUser) {
    return await super.create(user);
  }

  /**
   * Lists all users with a certain role
   * @param roles - The roles of user to list
   * @returns {Promise<User[]>} The found users
   */
  async listUsersWithRole(roles: string[]) {
    const json = await this.http.getJSON<ModelJson<User>[]>(`/api${this.baseURL}`, {searchParams: {role: roles.join("|")}});
    return json.map(item => new this.model(item));
  }
}