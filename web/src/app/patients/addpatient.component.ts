import { Component, OnInit, Input, Output, EventEmitter } from "@angular/core";
import { NewUser, User, UserService } from "../core/services/user.service";

/**
 * Component which contains the create patient form
 */
@Component({
  selector: "add-patient",
  templateUrl: "./addpatient.component.html"
})
export class AddPatientComponent implements OnInit {
  @Input() closeFn: () => void;
  @Output() onCreate: EventEmitter<User> = new EventEmitter<User>();

  /**
   * the new user data
   * @type {NewUser}
   */
  newPatient: NewUser = {
    username: "",
    fullName: "",
    role: "patient",
    email: "",
    birthdate: new Date(),
    gender: "male",
    phone: "",

    doctorIds: []
  };

  birthdate: {year: number, month: number, day: number};

  /**
   * List of doctors shown in the 'doctors' dropdown
   */
  doctors: User[];

  constructor(private userService: UserService) {

  }

  /**
   * Initialization Angular lifecycle hook for {@link OnInit}
   */
  async ngOnInit() {
    this.doctors = await this.userService.listUsersWithRole(["doctor"]);
  }

  /**
   * Close closes the modal
   */
  close() {
    this.closeFn();
  }

  /**
   * Adds a doctor to the new patient
   */
  addDoctor() {
    this.newPatient.doctorIds.push(0);
  }

  /**
   * Removes a doctor from the new dose
   * @param i - The index of the doctor to remove
   */
  removeDoctor(i: number) {
    this.newPatient.doctorIds = this.newPatient.doctorIds.filter((val, idx) => i != idx);
  }

  /**
   * Creates the new patient
   */
  async createPatient() {
    this.newPatient.birthdate = new Date(`${this.birthdate.year}-${this.birthdate.month}-${this.birthdate.day}`);
    this.newPatient.doctorIds = this.newPatient.doctorIds.map(v => +v);
    this.onCreate.emit(await this.userService.create(this.newPatient));
    this.closeFn();
  }
}