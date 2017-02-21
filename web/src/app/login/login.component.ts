import { Component } from '@angular/core';
import { Router } from "@angular/router";

import { Credentials, emptyCredentials, AuthService } from "../core/services/auth.service";

/**
 * Component that contains the login form
 */
@Component({
  templateUrl: "./login.component.html",
  selector: "login",
  styleUrls: ["./login.component.scss"]
})
export class LoginComponent {
  /**
   * Login credentials
   * @type {Credentials}
   */
  credentials: Credentials = emptyCredentials();

  /**
   * Eventual error that occurred while logging in
   * @type {string}
   */
  loginError: string = "";

  constructor(private authService: AuthService, private router: Router) {

  }

  /**
   * Event handler for a user login event
   */
  async login() {
    // Authenticate the user
    try {
      await this.authService.authenticate(this.credentials);
      this.router.navigate(["/home"]);
    } catch(e) {
      this.loginError = e.message;
    }
  }
}