import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { Subscription } from 'rxjs';

import { MedicationService, Medication } from "../core/services/medication.service";

/**
 * Component for a single medication
 */
@Component({
  selector: "medication",
  templateUrl: "./medication.component.html"
})
export class MedicationComponent implements OnInit, OnDestroy {
  medication: Medication;
  routeDataSubscription: Subscription;

  constructor(private medicationService: MedicationService, private route: ActivatedRoute) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  ngOnInit() {
    this.routeDataSubscription = this.route.data.subscribe((data: {medication: Medication}) => {
      this.medication = data.medication;
    });
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.routeDataSubscription && this.routeDataSubscription.unsubscribe();
  }
}