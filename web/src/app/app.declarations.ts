import { AppComponent } from "./app.component";
import { LoginComponent } from "./login/login.component";
import { HomeComponent } from "./home/home.component";
import { MedicationsComponent } from "./medications/medications.component";
import { MedicationComponent } from "./medications/medication.component";

/**
 * All components that need to be imported into the App module
 * @type {Array}
 */
export const appComponents = [
  AppComponent,
  LoginComponent,
  HomeComponent,
  MedicationsComponent, MedicationComponent
];