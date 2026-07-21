import {
  convertDateToDateObjectByTimezone,
  convertDateToTimestampByTimezone,
  getLocalTimezone,
  isStartOfDay,
} from '@/helpers/date/date'

import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array'

export const pbehaviorExceptionToForm = (exception = {}, timezone = getLocalTimezone()) => ({
  name: exception.name ?? '',
  description: exception.description ?? '',
  visible: !(exception.hidden ?? false),
  exdates: exception.exdates
    ? addKeyInEntities(
        exception.exdates.map(({ begin, end, type }) => ({
          begin: convertDateToDateObjectByTimezone(begin, timezone),
          end: convertDateToDateObjectByTimezone(isStartOfDay(end) ? end - 1 : end, timezone),
          type: { ...type },
        })),
      )
    : [],
  _id: exception._id,
})

export const pbehaviorExceptionImportToForm = (exceptionImport = {}) => ({
  name: exceptionImport.name ?? '',
  type: exceptionImport.type ?? '',
  file: exceptionImport.file,
})

export const formToPbehaviorException = (exceptionForm = {}, timezone = getLocalTimezone()) => {
  const { exdates, visible, ...form } = exceptionForm

  return {
    exdates: removeKeyFromEntities(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: convertDateToTimestampByTimezone(begin, timezone),
      end: convertDateToTimestampByTimezone(end, timezone),
    })),
    hidden: !visible,
    ...form,
  }
}
