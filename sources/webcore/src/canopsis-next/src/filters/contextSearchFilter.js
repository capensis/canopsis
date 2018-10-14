export default function (text) {
  return `{"$and":[{"$or":[{"name":{"$regex":"${text}","$options":"i"}},
      {"type":{"$regex":"${text}","$options":"i"}}]}]}`;
}

