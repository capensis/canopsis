import { exceptionsToRequest, exdatesToRequest } from '@/helpers/forms/planning-pbehavior';

/**
 * Convert pbehavior to timespan
 *
 * @param {Object} pbehavior
 * @param {Number} viewFrom
 * @param {Number} viewTo
 * @param {Boolean} byDate
 * @return {Object}}
 */
export const pbehaviorToTimespan = ({
  pbehavior,
  viewFrom,
  viewTo,
  byDate = false,
}) => ({
  rrule: pbehavior.rrule,
  start_at: pbehavior.tstart,
  end_at: pbehavior.tstop,
  view_from: viewFrom,
  view_to: viewTo,
  exdates: exdatesToRequest(pbehavior.exdates),
  exceptions: exceptionsToRequest(pbehavior.exceptions),
  by_date: byDate,
});
