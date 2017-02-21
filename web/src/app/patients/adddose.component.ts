import { Component, Input, OnInit, Output, EventEmitter} from '@angular/core';
import { DoseService, NewDose, Dose } from "../core/services/dose.service";
import { MedicationService, Medication } from "../core/services/medication.service";

/**
 * Component that contains the add dose form
 */
@Component({
  selector: "add-dose",
  templateUrl: "./adddose.component.html"
})
export class AddDoseComponent implements OnInit {
  @Input() closeFn: () => any;
  @Input() patientId: number;
  @Output() onCreate: EventEmitter<Dose> = new EventEmitter<Dose>();

  medications: Medication[] = [];

  newDose: NewDose = {
    title: "",
    description: "",
    dispenseAfter: {hour: 0,  minute: 0},
    dispenseBefore: {hour: 0, minute: 0},
    medications: []
  };

  constructor(private doseService: DoseService, private medicationService: MedicationService) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  async ngOnInit() {
    this.medications = await this.medicationService.list();
  }

  /**
   * Close closes the modal
   */
  close() {
    this.closeFn();
  }

  /**
   * Adds a medication to the new dose
   */
  addMedication() {
    this.newDose.medications.push({amount: 1, medicationId: 0})
  }

  /**
   * Removes a medication from the new dose
   * @param i - The of the dose index to remove
   */
  removeMedication(i) {
    this.newDose.medications = this.newDose.medications.filter((val, idx) => i != idx);
  }

  /**
   * Creates the new dose
   */
  async createDose() {
    this.newDose.medications.forEach(medication => {
      medication.medicationId = parseInt(medication.medicationId.toString());
    });

    this.onCreate.emit(await this.doseService.create(this.patientId, this.newDose));

    this.closeFn();
  }
}