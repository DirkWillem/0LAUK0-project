import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';

import { Dose, DoseService, DoseMedication } from "../core/services/dose.service";
import { MedicationService, Medication } from "../core/services/medication.service";
import { PartialModel } from "../core/model";
import { NgbTimeStruct } from "@ng-bootstrap/ng-bootstrap";

/**
 * Component that handles the updating of a dose
 */
@Component({
  selector: "update-dose",
  templateUrl: "updatedose.component.html"
})
export class UpdateDoseComponent implements OnInit {
  @Input() patientId: number;
  @Input() dose: Dose;
  @Input() closeFn: () => any;
  medications: Medication[] = [];

  dispenseAfter: NgbTimeStruct = {
    hour: 0,
    minute: 0,
    second: 0
  };

  dispenseBefore: NgbTimeStruct = {
    hour: 0,
    minute: 0,
    second: 0
  };

  @Output() onUpdate: EventEmitter<Dose> = new EventEmitter<Dose>();

  constructor(private doseService: DoseService, private medicationService: MedicationService) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  async ngOnInit() {
    // Make a copy of the dose so we aren't updating the actual model
    this.dose = new Dose(this.dose.toJSON());
    this.dispenseAfter = this.parseTime(this.dose.dispenseAfter);
    this.dispenseBefore = this.parseTime(this.dose.dispenseBefore);

    this.medications = await this.medicationService.list();
  }

  /**
   * Close closes the modal
   */
  close() {
    this.closeFn();
  }

  /**
   * Adds a medication to the dose
   */
  addMedication() {
    this.dose.medications.push(new DoseMedication({amount: 1, medication: { id: 0, title: "", description: ""}}))
  }

  /**
   * Removes a medication from the dose
   * @param i - The of the dose index to remove
   */
  removeMedication(i) {
    this.dose.medications = this.dose.medications.filter((val, idx) => i != idx);
  }

  /**
   * Sends the updated dose to the service
   */
  async updateDose() {
    this.dose.dispenseAfter = this.formatTime(this.dispenseAfter);
    this.dose.dispenseBefore = this.formatTime(this.dispenseBefore);

    this.onUpdate.emit(await this.doseService.update(this.patientId, this.dose.id, this.dose));

    this.closeFn();
  }

  /**
   * Transforms a time string into an NgbTimeStruct
   * @param time - The time string
   * @returns {NgbTimeStruct} the parsed struct
   */
  private parseTime(time: string): NgbTimeStruct {
    const parts = time.split(":").map(val => parseInt(val));
    return {
      hour: parts[0],
      minute: parts[1],
      second: parts[2]
    }
  }

  /**
   * Transforms an NgbTimeStruct into a time string
   * @param time - The time struct
   * @returns {string} The formatted string
   */
  private formatTime(time: NgbTimeStruct): string {
    const pad = (n: number) => n.toString().length == 1 ? `0${n.toString()}` : n.toString();

    return `${pad(time.hour)}:${pad(time.minute)}:${pad(time.second)}`;
  }
}