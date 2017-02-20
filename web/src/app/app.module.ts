import { NgModule } from '@angular/core';
import { BrowserModule } from "@angular/platform-browser";
import { FormsModule } from "@angular/forms";
import { HttpModule } from "@angular/http";

import { appComponents } from "./app.declarations";
import { AppComponent } from "./app.component";
import { MaterialModule } from "@angular/material";

/**
 * Main application module
 */
@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MaterialModule
  ],

  providers: [],
  declarations: [...appComponents],
  bootstrap: [AppComponent]
})
export class AppModule {

}