import { Component } from '@angular/core';

import { AuthService } from "../core/services/auth.service";

/**
 * Represents a single item in the top menu bar
 */
interface MenuItem {
  title: string;
  path: string;
  roles: string[];
}

/**
 * Component that contains the home screen of the application
 */
@Component({
  selector: "home",
  templateUrl: "./home.component.html"
})
export class HomeComponent {
  /**
   * Contains all menu items shown in the top menu bar
   * @type {MenuItem[]}
   */
  private menuItems: MenuItem[] = [
    { title: "Medications", path: "medications", roles: ["admin", "doctor", "pharmacist"] },
    { title: "Patients", path: "patients", roles: ["admin", "doctor"] }
  ];

  constructor(private authService: AuthService) {
    setTimeout(() => {
      console.log(this.visibleMenuItems);

    });
  }

  /**
   * Returns all menu items that are visible to the user based in its current role
   * @returns {MenuItem[]} The visible items
   */
  get visibleMenuItems(): MenuItem[] {
    return this.menuItems.filter(item => item.roles.includes(this.authService.session.role));
  }
}