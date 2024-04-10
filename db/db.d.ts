export interface Locker {
  id: string;
}

export interface Registration {
  locker: string;
  user: string;
  name: string;
  expiry: string;
  expiryEmailSent: string | null;
}

export interface User {
  email: string;
}

export interface DB {
  locker: Locker;
  registration: Registration;
  user: User;
}
