import { RouterModule, Routes } from '@angular/router';

import { AutoLoginGuard } from "./core/guards/autologin.guard";
import { AuthGuard } from "./core/guards/auth.guard";
import { MedicationsResolve, MedicationResolve } from "./core/resolves/medication.resolve";

import { LoginComponent } from "./login/login.component";
import { HomeComponent } from "./home/home.component";
import { MedicationsComponent } from "./medications/medications.component";
import { MedicationComponent } from "./medications/medication.component";
import { PatientsComponent } from "./patients/patients.component";
import { PatientsResolve, UserResolve } from "./core/resolves/user.resolve";
import { PatientComponent } from "./patients/patient.component";
import { DosesResolve } from "./core/resolves/dose.resolve";
import { DoseSummariesResolve } from "./core/resolves/dosesummary.resolve";
import { PRNMedicationsResolve } from "./core/resolves/prnmedication.resolve";

/**
 * Contains all routes of the application
 * @type {Routes}
 */
const routes = [
  {path: "", pathMatch: "full", redirectTo: "login"},
  {
    path: "login",
    component: LoginComponent,
    canActivate: [AutoLoginGuard]
  },
  {
    path: "home",
    component: HomeComponent,
    canActivate: [AuthGuard],
    children: [
      {
        path: "medications",
        component: MedicationsComponent,
        resolve: {
          medications: MedicationsResolve
        }
      },
      {
        path: "medications/:id",
        component: MedicationComponent,
        resolve: {
          medication: MedicationResolve
        }
      },
      {
        path: "patients",
        component: PatientsComponent,
        resolve: {
          patients: PatientsResolve
        }
      },
      {
        path: "patients/:id",
        component: PatientComponent,
        resolve: {
          patient: UserResolve,
          doses: DosesResolve,
          doseSummaries: DoseSummariesResolve,
          prnMedications: PRNMedicationsResolve
        }
      }
    ]
  }
];

/**
 * The application routing module
 * @type {ModuleWithProviders}
 */
export const appRoutingModule = RouterModule.forRoot(routes);