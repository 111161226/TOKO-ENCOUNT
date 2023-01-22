export const reorderArray = <T>(array: T[], from: number, to: number) => {
  const copyArray = array.slice()
  const item = copyArray[from]
  copyArray.splice(from, 1)
  copyArray.splice(to, 0, item)
  return copyArray
}
