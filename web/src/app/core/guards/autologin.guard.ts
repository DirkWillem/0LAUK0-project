import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

/**
 * Route guards that automatically redirects the user to the home screen if they are already logged in
 */
@Injectable()
export class AutoLoginGuard implements CanActivate {
  constructor(private router: Router) {

  }

  /**
   * Checks whether the current route may be activated
   * @returns {boolean} - Whether the route can be activated
   */
  canActivate(): boolean {
    if(window.sessionStorage.getItem("jwt")) {
      this.router.navigate(["/home"]);
      return true;
    }

    return true;
  }
}