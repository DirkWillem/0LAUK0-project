import { platformBrowserDynamic } from "@angular/platform-browser-dynamic";
import 'rxjs/add/operator/toPromise';

import { AppModule } from "./app/app.module";

// Run the app module
platformBrowserDynamic().bootstrapModule(AppModule);