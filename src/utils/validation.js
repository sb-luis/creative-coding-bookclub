export function validatePathsAreInScope(modifiedFiles, scope) {
  const invalidFiles = modifiedFiles.filter((file) => !file.startsWith(scope))

  if (invalidFiles.length > 0) {
    throw new Error(
      `Found ${invalidFiles.length} file/s outside the allowed member scope: '${scope}'`,
    )
  }
}

export function validateFilePath(filePath) {
  // Validate file path length
  const maxFilePathLength = 100
  const isValidLength = filePath.length <= maxFilePathLength
  if (!isValidLength) {
    throw new Error(
      `All file paths must be less than ${maxFilePathLength} characters. Looks like one of yours is ${filePath.length} characters long`,
    )
  }

  // Split file from rest of the path
  const fileNameFull = filePath.split('/').pop() ?? filePath

  // Split file name and extension
  let filename = fileNameFull
  let extension = null
  const lastDotIndex = fileNameFull.lastIndexOf('.')
  if (lastDotIndex !== -1) {
    filename = fileNameFull.substring(0, lastDotIndex)
    extension = fileNameFull.substring(lastDotIndex + 1)
  }

  // Validate file and folder names
  const filePathFragments = filePath.split('/') ?? [filePath]
  for (let i = 0; i < filePathFragments.length; i++) {
    let str = filePathFragments[i]
    if (str === '' && i === 0) continue
    if (i === filePathFragments.length - 1) str = filename
    const isKebabCase = /^[a-zA-Z0-9_-]+$/.test(str)
    if (!isKebabCase) {
      throw new Error(
        `Only alphanumeric characters, hyphens, and underscores are allowed for files and folder names.`,
      )
    }
  }

  // Validate extension is one of the allowed values
  const allowedExtensions = ['js', 'json', 'astro']
  const isValidExtension = allowedExtensions.includes(extension)
  if (!isValidExtension) {
    throw new Error(
      `Only the following file extensions are allowed: ${allowedExtensions.join(', ')}`,
    )
  }
}
