import { Injectable } from "@angular/core";
import { Resolve } from "@angular/router";

import { User, UserService } from "../services/user.service";
import { SingleResolve } from "./single.resolve";

@Injectable()
export class PatientsResolve implements Resolve<User[]> {
  constructor(private userService: UserService) {
  }

  async resolve(): Promise<User[]> {
    return await this.userService.listUsersWithRole(["patient"]);
  }
}

@Injectable()
export class UserResolve extends SingleResolve<User> {
  constructor(service: UserService) {
    super(service);
  }
}