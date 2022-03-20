import { names } from "./routes";
import router from "./router";

export function goTo(
  name: string,
  params: {
    replace: boolean;
  }
): void {
  router.push({
    name,
    replace: params.replace,
  });
}

export function goToLogin(replace?: boolean): void {
  goTo(names.login, {
    replace: replace ?? false,
  });
}

export function goToRegister(replace?: boolean): void {
  goTo(names.register, {
    replace: replace ?? false,
  });
}

export function goToHome(replace?: boolean): void {
  goTo(names.home, {
    replace: replace ?? false,
  });
}

export function goToProfile(replace?: boolean): void {
  goTo(names.profile, {
    replace: replace ?? false,
  });
}
