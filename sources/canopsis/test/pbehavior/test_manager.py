class TestBla(TestCase):

    def test_check_active_pbehavior(self):
        now = int(time.mktime(datetime.utcnow().timetuple()))
        hour = 3600

        # tstart < now < tstop
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour,
            now + hour,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertTrue(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # tstart is one hour ahead from now
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now + hour * 1,
            now + hour * 2,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # tstart is two hour behind from now, tstop one hour
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 2,
            now - hour * 1,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # no rrule, tstart and tstop in the past
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 2,
            now - hour * 1,
            '',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_n_rrule))

        # no rrule, now between tstart and tstop
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 1,
            now + hour * 1,
            '',
            'test'
        ).to_dict()

        self.assertTrue(self.pbm.check_active_pbehavior(now, pb_n_rrule))

        # no rrule, tstart and tstop in the future
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now + hour * 1,
            now + hour * 2,
            '',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_n_rrule))

    def test_check_active_pbehavior_2(self):
        timestamps = []

        timestamps.append((False, 1529154801-24*3600))  # Vendredi 15 Juin 2018 15h13
        timestamps.append((True, 1529154801-24*3600+5*3600))  # Vendredi 15 Juin 2018 20h13
        timestamps.append((True, 1529154801))  # Samedi 16 Juin 2018 15h13
        timestamps.append((True, 1529290800))  # Lundi 18 Juin 2018 05h00

        timestamps.append((False, 1529308800))  # Lundi 18 Juin 2018 10h00
        timestamps.append((False, 1529308800+7*24*3600))
        timestamps.append((False, 1529308800+7*24*3600*2))
        timestamps.append((False, 1529308800+7*24*3600*3))
        timestamps.append((False, 1529308800+7*24*3600*4))
        timestamps.append((False, 1529308800+7*24*3600*5))

        timestamps.append((True, 1529740800))  # Samedi 23 Juin 2018 10h00
        timestamps.append((True, 1529740800+7*24*3600)) # +7j
        timestamps.append((True, 1529740800+7*24*3600*2)) # ...
        timestamps.append((True, 1529740800+7*24*3600*3))
        timestamps.append((True, 1529740800+7*24*3600*4))
        timestamps.append((True, 1529740800+7*24*3600*5))

        pbehavior = {
            "rrule": "FREQ=WEEKLY;BYDAY=FR",
            "tstart": 1529085600,
            "tstop": 1529294400,
        }

        for i, ts in enumerate(timestamps):
            res = PBehaviorManager.check_active_pbehavior(ts[1], pbehavior)
            self.assertEqual(res, ts[0])

    def test_get_active_intervals(self):
        day = 24 * 3600
        tstart = 1530288000  # 2018/06/29 18:00:00
        tstop = tstart + 3600

        pbehavior = {
            'rrule': 'FREQ=DAILY',
            'tstart': tstart,
            'tstop': tstop
        }

        # after = tstart
        expected_intervals = [
            (tstart, tstop),
            (tstart + day, tstop + day),
            (tstart + 2 * day, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart, tstart + 5 * day, pbehavior))
        self.assertEqual(intervals, expected_intervals)

        # after < tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart - 3 * day, tstart + 5 * day, pbehavior))
        self.assertEqual(intervals, expected_intervals)

        # after > tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart + 2 * day, tstart + 5 * day, pbehavior))
        expected_intervals = [
            (tstart + 2 * day, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        self.assertEqual(intervals, expected_intervals)

        intervals = list(PBehaviorManager.get_active_intervals(
            tstart + 2 * day + 1800, tstart + 5 * day, pbehavior))
        expected_intervals = [
            (tstart + 2 * day + 1800, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        self.assertEqual(intervals, expected_intervals)

        # before < tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart - 3 * day, tstart - 2 * day, pbehavior))
        expected_intervals = []
        self.assertEqual(intervals, expected_intervals)

    def test_get_intervals_with_pbehaviors(self):
        day = 24 * 3600

        tstart1 = 1530288000  # 2018/06/29 18:00:00
        tstop1 = tstart1 + 3600

        tstart2 = tstart1 + 1800
        tstop2 = tstop1 + 1800

        pbehavior1 = deepcopy(self.pbehavior)
        pbehavior2 = deepcopy(self.pbehavior)
        pbehavior1.update({
            'eids': [1],
            'rrule': 'FREQ=DAILY',
            'tstart': tstart1,
            'tstop': tstop1
        })
        pbehavior2.update({
            'eids': [1],
            'rrule': 'FREQ=DAILY',
            'tstart': tstart2,
            'tstop': tstop2
        })

        self.pbm.pb_storage.put_elements(
            elements=(pbehavior1, pbehavior2)
        )

        expected_intervals = [
            (tstart1, tstart1, False),
            (tstart1, tstop2, True),
            (tstop2, tstart1 + day, False),
            (tstart1 + day, tstop2 + day, True),
            (tstop2 + day, tstart1 + 2 * day, False),
            (tstart1 + 2 * day, tstop2 + 2 * day, True),
            (tstop2 + 2 * day, tstart1 + 3 * day, False),
            (tstart1 + 3 * day, tstop2 + 3 * day, True),
            (tstop2 + 3 * day, tstart1 + 4 * day, False),
            (tstart1 + 4 * day, tstop2 + 4 * day, True),
            (tstop2 + 4 * day, tstart1 + 5 * day, False),
        ]
        intervals = list(self.pbm.get_intervals_with_pbehaviors(
            tstart1, tstart1 + 5 * day, 1))
        self.assertEqual(intervals, expected_intervals)

        # Entity without pbehaviors
        expected_intervals = [
            (tstart1, tstart1 + 5 * day, False),
        ]
        intervals = list(self.pbm.get_intervals_with_pbehaviors(
            tstart1, tstart1 + 5 * day, 2))
        self.assertEqual(intervals, expected_intervals)
