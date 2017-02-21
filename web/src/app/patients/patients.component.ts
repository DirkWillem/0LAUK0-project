import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { Subscription } from 'rxjs';

import { User, UserService } from "../core/services/user.service";

/**
 * Component containing a list of patients
 */
@Component({
  selector: "patients",
  templateUrl: "./patients.component.html"
})
export class PatientsComponent implements OnInit, OnDestroy {
  patients: User[] = [];
  routeDataSubscription: Subscription;

  constructor(private userService: UserService,
              private route: ActivatedRoute,
              private router: Router) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  ngOnInit() {
    this.routeDataSubscription = this.route.data.subscribe((data: {patients: User[]}) => {
      this.patients = data.patients;
    });
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.routeDataSubscription && this.routeDataSubscription.unsubscribe();
  }

}