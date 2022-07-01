import { exceptionsToRequest, exdatesToRequest } from '@/helpers/forms/planning-pbehavior';

/**
 * Convert pbehavior to timespan
 *
 * @param {Object} pbehavior
 * @param {number} from
 * @param {number} to
 * @param {boolean} byDate
 * @return {Object}}
 */
export const pbehaviorToTimespanRequest = ({
  pbehavior,
  from,
  to,
  byDate = false,
}) => {
  const tstartBeforeFrom = pbehavior.tstart < from;
  const tstopAfterFrom = !pbehavior.tstop || (pbehavior.tstop > from);

  const tstartBeforeTo = pbehavior.tstart < to;
  const tstopAfterTo = pbehavior.tstop && (pbehavior.tstop > to);

  const viewFrom = (tstartBeforeFrom && tstopAfterFrom) ? pbehavior.tstart : from;
  const viewTo = (tstartBeforeTo && tstopAfterTo) ? pbehavior.tstop : to;

  const request = {
    rrule: pbehavior.rrule,
    start_at: pbehavior.tstart,
    type: pbehavior.type._id,
    view_from: viewFrom,
    view_to: viewTo,
    exdates: exdatesToRequest(pbehavior.exdates),
    exceptions: exceptionsToRequest(pbehavior.exceptions),
    by_date: byDate,
  };

  if (pbehavior.tstop) {
    request.end_at = pbehavior.tstop;
  }

  return request;
};
