import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { Subscription } from 'rxjs';

import { User, UserService } from "../core/services/user.service";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";

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
              private router: Router,
              private modalService: NgbModal) {

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

  /**
   * Opens the add-patient modal
   * @param content - Reference to the modal contents
   */
  openAddPatientModal(contents) {
    this.modalService.open(contents, {size: "lg"});
  }

  /**
   * Adds a patient to the list of patients if it wasn't already added by the dispatcher
   * @param patient - The created patient
   */
  patientCreated(patient: User) {
    if(!this.patients.some(p => p.id == patient.id)) {
      this.patients.push(patient);
    }
    this.router.navigate([patient.id], {relativeTo: this.route});
  }

}