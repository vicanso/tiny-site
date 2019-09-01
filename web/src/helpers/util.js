import { sha256 } from "./crypto";

const hash = "tiny-site";
const suRole = "su";
const adminRole = "role";

export function setBeginOfDay(date) {
  return date
    .clone()
    .hours(0)
    .minutes(0)
    .seconds(0)
    .milliseconds(0);
}

export function setEndOfDay(date) {
  return date
    .clone()
    .hours(23)
    .minutes(59)
    .seconds(59)
    .milliseconds(999);
}

export function generatePassword(pass) {
  return sha256(pass + hash);
}

// includeRole 判断是否包括角色
function includeRole(roles, checkRoles) {
  let found = false;
  if (!roles || !checkRoles) {
    return found;
  }
  roles.forEach(item => {
    if (found) {
      return;
    }
    checkRoles.forEach(checkItem => {
      if (item === checkItem) {
        found = true;
      }
    });
  });
  return found;
}

// isAdminUser 判断是否admin
export function isAdminUser(roles) {
  return includeRole(roles, [suRole, adminRole]);
}
