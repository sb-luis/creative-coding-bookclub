import { glob } from "glob"
import { copyFile } from "fs";

const read = async () => {
  glob("./src/members/**/*.[jt]s")
    .then(async res => {
      return await Promise.all(
        res.map(f => {
          return new Promise((resolve, rej) => {
            const fileParts = f.split('/');
            console.log(f)
            copyFile(f, "dist/source-"+fileParts[fileParts.length - 2]+"-"+fileParts[fileParts.length - 1], (err, res) => {
              resolve(res);
            })
          })
      }))
    })
}
read();