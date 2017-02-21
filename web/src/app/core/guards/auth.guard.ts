import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

/**
 * Route guard that checks whether the user is logged in, or redirects them back to the login page
 */
@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private router: Router) {

  }

  /**
   * Checks whether the current route may be activated
   * @returns {boolean} - Whether the route can be activated
   */
  canActivate(): boolean {
    if(window.sessionStorage.getItem("jwt")) {
      return true;
    }

    this.router.navigate(["/login"]);
    return false;
  }
}