import { AppComponent } from "./app.component";
import { LoginComponent } from "./login/login.component";
import { HomeComponent } from "./home/home.component";
import { MedicationsComponent } from "./medications/medications.component";
import { MedicationComponent } from "./medications/medication.component";
import { PatientsComponent } from "./patients/patients.component";
import { PatientComponent } from "./patients/patient.component";
import { AddDoseComponent } from "./patients/adddose.component";
import { UpdateDoseComponent } from "./patients/updatedose.component";
import { DoseSummaryComponent } from "./patients/dosesummary.component";
import { AddPatientComponent } from "./patients/addpatient.component";
import { AddPRNMedicationComponent } from "./patients/addprnmedication.component";

/**
 * All components that need to be imported into the App module
 * @type {Array}
 */
export const appComponents = [
  AppComponent,
  LoginComponent,
  HomeComponent,
  MedicationsComponent, MedicationComponent,
  PatientsComponent, PatientComponent, AddPatientComponent,
    AddDoseComponent, UpdateDoseComponent, DoseSummaryComponent,
    AddPRNMedicationComponent
];