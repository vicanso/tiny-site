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

// getQueryParams 获取query string的参数
export function getQueryParams(searcValue, key) {
  if (!searcValue) {
    return "";
  }
  if (searcValue[0] === "?") {
    searcValue = searcValue.substring(1);
  }
  const arr = searcValue.split("&");
  let result = "";
  arr.forEach(item => {
    if (result) {
      return;
    }
    const tmpArr = item.split("=");
    if (tmpArr[0] === key) {
      result = tmpArr[1];
    }
  });
  return result;
}

// copy 将内容复制至粘贴板
export function copy(value, parent) {
  if (!document.execCommand) {
    return new Error("该浏览器不支持复制功能");
  }
  const input = document.createElement("input");
  input.value = value;
  if (parent) {
    parent.appendChild(input);
  } else {
    document.body.appendChild(input);
  }
  input.focus();
  input.select();
  document.execCommand("Copy", false, null);
  input.remove();
}
