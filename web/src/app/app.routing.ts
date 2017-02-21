import { RouterModule, Routes } from '@angular/router';

import { AutoLoginGuard } from "./core/guards/autologin.guard";
import { AuthGuard } from "./core/guards/auth.guard";
import { MedicationsResolve, MedicationResolve } from "./core/resolves/medication.resolve";

import { LoginComponent } from "./login/login.component";
import { HomeComponent } from "./home/home.component";
import { MedicationsComponent } from "./medications/medications.component";
import { MedicationComponent } from "./medications/medication.component";

/**
 * Contains all routes of the application
 * @type {Routes}
 */
const routes: Routes = [
  { path: "", pathMatch: "full", redirectTo: "login" },
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
      }
    ]
  }
];

/**
 * The application routing module
 * @type {ModuleWithProviders}
 */
export const appRoutingModule = RouterModule.forRoot(routes);