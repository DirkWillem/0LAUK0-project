import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import * as decodeJwt from 'jwt-decode';

/**
 * Interface containing user credentials
 */
export interface Credentials {
  username: string;
  password: string;
}

/**
 * Returns an empty set of credentials
 * @returns {Credentials} - The empty credentials
 */
export function emptyCredentials(): Credentials {
  return {
    username: "",
    password: ""
  };
}

/**
 * Response that is returned from an authentication request
 */
interface TokenResponse {
  token: string;
}

/**
 * Contains session information
 */
export interface Session {
  userId: number;
  username: string;
  fullName: string;
  role: string;
  email: string;
}

/**
 * Service that handles authentication
 */
@Injectable()
export class AuthService {
  constructor(private http: Http) { }

  /**
   * Authenticates the user with the given credentials
   * @param credentials - Credentials to authenticate the user with
   * @returns {Promise<void>} - Promise that resolves on a successful login, or rejects on a failed login
   */
  async authenticate(credentials: Credentials): Promise<void> {
    try {
      const response = await this.http.post("/api/authenticate", credentials).toPromise();
      const json: TokenResponse = response.json();

      window.sessionStorage.setItem("jwt", json.token);
    } catch(e) {
      if(e instanceof Response) {
        throw new Error(e.json().message);
      }
      throw e;
    }
  }

  /**
   * Returns the current session
   */
  get session(): Session {
    const jwt = sessionStorage.getItem("jwt");
    if(jwt) {
      return decodeJwt(sessionStorage.getItem("jwt"));
    } else {
      return null;
    }
  }
}