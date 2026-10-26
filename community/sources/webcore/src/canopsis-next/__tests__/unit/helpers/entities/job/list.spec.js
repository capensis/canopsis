import {
  prepareJobsStopRequest,
  splitJobsBulkResponse,
} from '@/helpers/entities/job/list';

describe('helpers/entities/job/list', () => {
  describe('prepareJobsStopRequest', () => {
    const jobs = [
      { _id: 'first', pending_jobs: 2 },
      { _id: 'second', pending_jobs: 0 },
    ];

    it('adds a stop_pending_jobs decision for remediation and webhook jobs', () => {
      expect(prepareJobsStopRequest(jobs, { first: true })).toEqual([
        { _id: 'first', stop_pending_jobs: true },
        { _id: 'second', stop_pending_jobs: false },
      ]);
    });

    it('does not add stop_pending_jobs for auth token jobs', () => {
      expect(prepareJobsStopRequest(jobs, {}, false)).toEqual([
        { _id: 'first' },
        { _id: 'second' },
      ]);
    });
  });

  describe('splitJobsBulkResponse', () => {
    it('separates successful and failed bulk response items', () => {
      const success = { id: 'first', status: 200 };
      const failed = { item: { _id: 'second' }, error: 'not found', status: 404 };

      expect(splitJobsBulkResponse([success, failed])).toEqual({
        successItems: [success],
        failedItems: [failed],
      });
    });
  });
});
