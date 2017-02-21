import { Injectable } from '@angular/core';
import { Http, Headers, Response, RequestOptions, ResponseContentType, URLSearchParams } from '@angular/http';

/**
 * Returns whether a value with possible types URLSearchParams or {[key: string]: any} is of type URLSearchParams
 * @param s - The value to determine the type of
 * @returns {boolean} Whether the value is of type URLSearchParams
 */
function isURLSearchParams(s: URLSearchParams | {[key: string]: any}): s is URLSearchParams {
  return s instanceof URLSearchParams
}

/**
 * Contains the base options shared across all request types
 */
export interface BaseRequestOptions {
  searchParams?: URLSearchParams | {[key: string]: any};
  headers?: {[key: string]: string};
}


/**
 * Service for handling HTTP calls that require a session
 */
@Injectable()
export class AuthHttp {
  constructor(private http: Http) { }

  /**
   * Executes a POST request to a JSON API endpoint
   * @param url - The URL of the API endpoint
   * @param data - The data to be sent
   * @param options - Additional options for the request (optional)
   * @returns {Promise<T>} The parsed JSON returned from the API
   */
  async postJSON<T>(url: string, data: any, options?: BaseRequestOptions): Promise<T> {
    const response = await this.http.post(url, data, this.createRequestOptions(options || {})).toPromise();
    return response.json() as T;
  }

  /**
   * Executes a GET request to a JSON API endpoint
   * @param url - The URL of the API endpoint
   * @param options - Additional options for the request (optional)
   * @returns {Promise<T>} The parsed JSON returned from the API
   */
  async getJSON<T>(url: string, options?: BaseRequestOptions): Promise<T> {
    const response = await this.http.get(url, this.createRequestOptions(options || {})).toPromise();
    return response.json() as T;
  }

  /**
   * Returns the HTTP request options for a given list of base request options
   * @param options - The request options
   * @returns {RequestOptions} The created request options
   */
  private createRequestOptions(options: BaseRequestOptions): RequestOptions {
    // Create the base requests object
    const opts: any = {
      headers: this.createHeaders(options)
    };

    // Add search params if necessary
    if(options.searchParams) {
      if(isURLSearchParams(options.searchParams)) {
        opts.search = options.searchParams;
      } else {
        const sp = new URLSearchParams();
        for(let param in options.searchParams) {
          sp.append(param, options.searchParams[param]);
        }
        opts.search = sp;
      }
    }

    return new RequestOptions(opts);
  }

  /**
   * Returns the default request headers and adds any headers set on the request options
   * @param options - The options for the request to create the headers for
   * @returns {Headers} The HTTP headers
   */
  private createHeaders(options: BaseRequestOptions): Headers {
    return new Headers(Object.assign({}, options.headers, {"X-JWT": sessionStorage.getItem("jwt")}));
  }
}