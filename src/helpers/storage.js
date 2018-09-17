const listFilePageSizeKey = "list-file-page-size";
export function saveListFilePageSize(size) {
  if (localStorage) {
    localStorage.setItem(listFilePageSizeKey, size);
  }
  return;
}

export function getListFilePageSize() {
  let pageSize = 10;
  if (localStorage) {
    const v = localStorage.getItem(listFilePageSizeKey);
    if (v) {
      pageSize = Number.parseInt(v, 10);
    }
  }
  return pageSize;
}
