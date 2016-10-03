#!/usr/bin/env python
# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2015 Jonathan Labéjof <jonathan.labejof@gmail.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
# --------------------------------------------------------------------

from unittest import main

from time import sleep

from inspect import getmembers

from b3j0f.utils.ut import UTCase

from six.moves import range

from ..core import Annotation, StopPropagation, RoutineAnnotation


class TestAnnotation(Annotation):
    """Annotation for inheritance tests."""


class AnnotationTest(UTCase):
    """UT class which creates an annotation and delete it at the end."""

    def setUp(self):
        """Create a new annotation."""

        self.annotation = Annotation()

    def tearDown(self):
        """Delete self.annotation."""

        self.annotation.__del__()
        del self.annotation


class MemoryTest(AnnotationTest):
    """Test InMemory Annotations."""

    def setUp(self):

        super(MemoryTest, self).setUp()

        Annotation.free_memory()

    def tearDown(self):

        super(MemoryTest, self).tearDown()

        Annotation.free_memory()

    def test_not_in_memory(self):
        """Test annotation not in memory
        """

        annotations = Annotation.get_memory_annotations()
        self.assertFalse(annotations)

    def test_in_memory(self):
        """Test if annotation is in memory
        """

        self.annotation = Annotation(in_memory=True)
        annotations = Annotation.get_memory_annotations()
        self.assertEqual(annotations, set((self.annotation,)))

    def test_set_in_memory(self):
        """Test set in_memory
        """

        self.annotation.in_memory = True
        annotations = Annotation.get_memory_annotations()
        self.assertEqual(annotations, set((self.annotation,)))

    def test_set_not_in_memory(self):
        """Test set not in_memory
        """

        self.annotation.in_memory = True
        annotations = Annotation.get_memory_annotations()
        self.assertEqual(annotations, set((self.annotation,)))
        self.annotation.in_memory = False
        annotations = Annotation.get_memory_annotations()
        self.assertFalse(annotations)

    def test_get(self):

        testAnnotation = TestAnnotation(in_memory=True)
        self.annotation.in_memory = True
        annotations = Annotation.get_memory_annotations()
        self.assertEqual(annotations, set((testAnnotation, self.annotation)))
        testAnnotation.__del__()

    def test_get_inheritance(self):

        testAnnotation = TestAnnotation(in_memory=True)
        self.annotation.in_memory = True
        annotations = TestAnnotation.get_memory_annotations()
        self.assertEqual(annotations, set((testAnnotation,)))
        testAnnotation.__del__()

    def test_exclude(self):

        testAnnotation = TestAnnotation(in_memory=True)
        self.annotation.in_memory = True
        annotations = Annotation.get_memory_annotations(exclude=TestAnnotation)
        self.assertEqual(annotations, set((self.annotation,)))
        annotations = Annotation.get_memory_annotations(exclude=Annotation)
        self.assertFalse(annotations)
        testAnnotation.__del__()


class DeleteTest(AnnotationTest):
    """Test annotation deletion."""

    def test_one_annotation(self):
        """Test to delete one bound annotation."""

        self.annotation(self)

        annotations = Annotation.get_annotations(self)

        self.assertEqual(len(annotations), 1)

        self.annotation.__del__()

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

    def test_two_annotations(self):
        """Test to delete one annotation bound twice on the same element."""

        self.annotation(self)
        self.annotation(self)

        annotations = Annotation.get_annotations(self)

        self.assertEqual(len(annotations), 2)

        self.annotation.__del__()

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

    def test_two_annotations_wit_two_objects(self):
        """Test to delete one annotations bound to two elements."""

        self.annotation(self)
        self.annotation(DeleteTest)

        annotations = Annotation.get_annotations(self)

        self.assertEqual(len(annotations), 2)

        annotations = Annotation.get_annotations(DeleteTest)

        self.assertEqual(len(annotations), 1)

        self.annotation.__del__()

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

        annotations = Annotation.get_annotations(DeleteTest)

        self.assertFalse(annotations)


class RemoveTest(AnnotationTest):
    """Test remove class method."""

    def setUp(self):

        super(RemoveTest, self).setUp()

        self.test_annotation = TestAnnotation()

        self.annotation(self)
        self.test_annotation(self)

    def tearDown(self):

        del self.annotation
        del self.test_annotation

    def test(self):
        """Test simple remove."""

        Annotation.remove(self)

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

    def test_inheritance(self):

        TestAnnotation.remove(self)

        annotations = Annotation.get_annotations(self)

        self.assertEqual(annotations[0], self.annotation)

    def test_exclude(self):

        Annotation.remove(self, exclude=TestAnnotation)

        annotations = Annotation.get_annotations(self)

        self.assertEqual(annotations[0], self.test_annotation)


class OnBindTargetTest(AnnotationTest):
    """Test on_bind_target handler."""

    def setUp(self):

        super(OnBindTargetTest, self).setUp()

        self.count = 0
        self.annotation._on_bind_target = self.on_bind_target

    def on_bind_target(self, annotation, target, ctx, *args, **kwargs):
        self.count += 1

    def test_one(self):

        self.annotation(self)

        self.assertEqual(self.count, 1)

    def test_many(self):

        self.annotation(self)
        self.annotation(self)

        self.assertEqual(self.count, 2)

    def test_many_elts(self):

        self.annotation(self)
        self.annotation(OnBindTargetTest)
        self.annotation(OnBindTargetTest)

        self.assertEqual(self.count, 3)


class TargetsTest(AnnotationTest):
    """Test targets attribute."""

    def test_none(self):

        self.assertFalse(self.annotation.targets)

    def test_one(self):

        self.annotation(self)

        self.assertIn(self, self.annotation.targets)

    def test_many(self):

        self.annotation(self)
        self.annotation(self)

        self.assertIn(self, self.annotation.targets)

    def test_many_many(self):

        self.annotation(self)
        self.annotation(TargetsTest)

        self.assertIn(self, self.annotation.targets)
        self.assertIn(TargetsTest, self.annotation.targets)


class TTLTest(AnnotationTest):
    """Test ttl."""

    def setUp(self):

        super(TTLTest, self).setUp()

        self.ttl = 0.1

    def test_def(self):
        """Test ttl at definition."""

        self.annotation = Annotation(ttl=self.ttl)
        self.annotation(self)

        annotations = Annotation.get_annotations(self)

        self.assertTrue(annotations)

        sleep(2 * self.ttl)

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

    def test_run(self):
        """Test to set ttl at runtime."""

        self.annotation(self)

        annotations = Annotation.get_annotations(self)

        self.assertTrue(annotations)

        self.annotation.ttl = self.ttl

        sleep(2 * self.ttl)

        annotations = Annotation.get_annotations(self)

        self.assertFalse(annotations)

    def test_run_run(self):
        """Test to change ttl after changing it a first time."""

        self.annotation(self)

        annotations = Annotation.get_annotations(self)

        self.assertTrue(annotations)

        self.annotation.ttl = 5

        self.assertLess(self.annotation.ttl, 5)

        self.annotation.ttl = 10

        self.assertGreater(self.annotation.ttl, 5)

        self.annotation.ttl = None

        self.assertIsNone(self.annotation.ttl)


class GetAnnotationsTest(AnnotationTest):
    """Test to annotate elements."""

    def test_None(self):
        """Test to annotate None."""

        self.annotation(None)

        annotations = Annotation.get_annotations(None)

        self.assertTrue(annotations)

    def test_number(self):
        """Test to annotate a number."""

        self.annotation(1)

        annotations = Annotation.get_annotations(1)

        self.assertTrue(annotations)

    def test_function(self):
        """Test to annotate a function."""

        def test():
            """test function."""

        self.annotation(test)

        annotations = Annotation.get_annotations(test)

        self.assertTrue(annotations)

    def test_builtin(self):
        """Test to annotate a builtin element."""

        self.annotation(range)

        annotations = Annotation.get_annotations(range)

        self.assertTrue(annotations)

    def test_class(self):
        """Test to annotate a class."""

        class Test(object):
            """Test class."""

        self.annotation(Test)

        annotations = Annotation.get_annotations(Test)

        self.assertTrue(annotations)

    def test_dclass(self):
        """Test to annotate directly a class."""

        @Annotation()
        class Test(object):
            """Test class."""

        annotations = Annotation.get_annotations(Test)

        self.assertTrue(annotations)

    def test_namespace(self):
        """Test to annotate a namespace."""

        class Test:
            """Test namespace."""

        self.annotation(Test)

        annotations = Annotation.get_annotations(Test)

        self.assertTrue(annotations)

    def test_method(self):
        """Test to annotate a method."""

        class Test:
            """Test class."""
            def test(self):
                """Test method."""

        self.annotation(Test.test, ctx=Test)

        annotations = Annotation.get_annotations(Test.test, ctx=Test)

        self.assertTrue(annotations)

    def test_dmethod(self):
        """Test to annotate directly a method."""

        class Test:
            """Test class."""
            @Annotation()
            def test(self):
                """Test method."""

        annotations = Annotation.get_annotations(Test.test, ctx=Test)

        self.assertTrue(annotations)

    def test_boundmethod(self):
        """Test to annotate a bound method."""

        class Test:
            """Test class."""
            def test(self):
                """Test method."""

        test = Test()

        self.annotation(test.test)

        annotations = Annotation.get_annotations(test.test)

        self.assertTrue(annotations)

    def test_instance(self):
        """Test to annotate an instance."""

        class Test:
            """Test class."""

        test = Test()

        self.annotation(test)

        annotations = Annotation.get_annotations(test)

        self.assertTrue(annotations)

    def test_module(self):
        """Test to annotate a module."""

        import sys

        self.annotation(sys)

        annotations = Annotation.get_annotations(sys)

        self.assertTrue(annotations)

    def test_dconstructor(self):
        """Test to annotate directly a constructor."""

        class Test(object):

            @Annotation()
            def __init__(self):
                pass

        #self.annotation(Test.__init__, ctx=Test)

        annotations = Annotation.get_annotations(Test.__init__, ctx=Test)

        self.assertTrue(annotations)

    def test_constructor(self):
        """Test to annotate a constructor."""

        class Test(object):

            def __init__(self):
                pass

        self.annotation(Test.__init__, ctx=Test)

        annotations = Annotation.get_annotations(Test.__init__, ctx=Test)

        self.assertTrue(annotations)

    def test_boundconstructor(self):
        """Test to annotate a bound constructor."""

        class Test(object):

            def __init__(self):
                pass

        test = Test()

        self.annotation(test.__init__, ctx=test)

        annotations = Annotation.get_annotations(test.__init__, ctx=test)

        self.assertTrue(annotations)

    def test_dboundconstructor(self):
        """Test to annotate directly a bound constructor."""

        class Test(object):
            @Annotation()
            def __init__(self):
                pass

        test = Test()

        annotations = Annotation.get_annotations(test.__init__, ctx=test)

        self.assertTrue(annotations)

    def test_depthsearch(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(Test0, maxdepth=5)

        self.assertEqual(len(_annotations), 5)
        self.assertNotIn(annotation3, _annotations)

    def test_depthsearchnotfollowannotated(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(
            Test0, maxdepth=5, followannotated=False
        )

        self.assertEqual(len(_annotations), 6)
        self.assertIn(annotation3, _annotations)

    def test_depthsearch_mindepth(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(
            Test0, mindepth=1, maxdepth=5
        )

        self.assertEqual(len(_annotations), 4)
        self.assertNotIn(annotation3, _annotations)

    def test_depthsearchnotfollowannotated_mindepth(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(
            Test0, mindepth=1, maxdepth=5, followannotated=False
        )

        self.assertEqual(len(_annotations), 5)
        self.assertIn(annotation3, _annotations)

    def test_depthsearch_public(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class _Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(Test0, maxdepth=5)

        self.assertEqual(len(_annotations), 2)
        self.assertNotIn(annotation2, _annotations)
        self.assertNotIn(annotation3, _annotations)

    def test_depthsearch_notpublic(self):

        annotations = [Annotation() for _ in range(4)]
        annotation0 = annotations[0]
        annotation1 = annotations[1]
        annotation2 = annotations[2]
        annotation3 = annotations[3]

        @annotation0
        class Test0(object):

            @annotation0
            @annotation1
            class _Test1(object):

                class Test(object):

                    @annotation3
                    def test(self):
                        pass

                @annotation2
                class Test2(object):

                    pass

            @annotation1
            def test1(self):
                pass

        _annotations = Annotation.get_annotations(
            Test0, maxdepth=5, public=False
        )

        self.assertEqual(len(_annotations), 5)
        self.assertNotIn(annotation3, _annotations)


class GetLocalAnnotationsTest(AnnotationTest):
    """Test get local annotatations."""

    def test_None(self):
        """Test to annotate None.
        """

        self.annotation(None)

        annotations = Annotation.get_local_annotations(None)

        self.assertTrue(annotations)

    def test_number(self):
        """Test to annotate a number.
        """

        self.annotation(1)

        annotations = Annotation.get_local_annotations(1)

        self.assertTrue(annotations)

    def test_function(self):
        """Test to annotate a function.
        """

        def test():
            pass

        self.annotation(test)

        annotations = Annotation.get_local_annotations(test)

        self.assertTrue(annotations)

    def test_builtin(self):
        """Test to annotate a builtin element.
        """

        self.annotation(range)

        annotations = Annotation.get_local_annotations(range)

        self.assertTrue(annotations)

    def test_class(self):
        """Test to annotate a class.
        """

        class Test(object):
            pass

        self.annotation(Test)

        annotations = Annotation.get_local_annotations(Test)

        self.assertTrue(annotations)

    def test_namespace(self):
        """Test to annotate a namespace.
        """

        class Test:
            pass

        self.annotation(Test)

        annotations = Annotation.get_local_annotations(Test)

        self.assertTrue(annotations)

    def test_method(self):
        """Test to annotate a method.
        """

        class Test:
            def test(self):
                pass

        self.annotation(Test.test)

        annotations = Annotation.get_local_annotations(Test.test)

        self.assertTrue(annotations)

    def test_boundmethod(self):
        """Test to annotate a bound method.
        """

        class Test:
            def test(self):
                pass

        test = Test()

        self.annotation(test.test)

        annotations = Annotation.get_local_annotations(test.test)

        self.assertTrue(annotations)

    def test_instance(self):
        """Test to annotate an instance.
        """

        class Test:
                pass

        test = Test()

        self.annotation(test)

        annotations = Annotation.get_local_annotations(test)

        self.assertTrue(annotations)

    def test_module(self):
        """Test to annotate a module.
        """

        import sys

        self.annotation(sys)

        annotations = Annotation.get_local_annotations(sys)

        self.assertTrue(annotations)

    def test_dconstructor(self):
        """Test to annotate directly a constructor.
        """

        class Test:
            @Annotation()
            def __init__(self):
                pass

        annotations = Annotation.get_local_annotations(Test.__init__)

        self.assertTrue(annotations)

    def test_constructor(self):
        """Test to annotate a constructor.
        """

        class Test:

            def __init__(self):
                pass

        self.annotation(Test.__init__)

        annotations = Annotation.get_local_annotations(Test.__init__)

        self.assertTrue(annotations)

    def test_boundconstructor(self):
        """Test to annotate a bound constructor.
        """

        class Test:
            def __init__(self):
                pass

        test = Test()

        self.annotation(test.__init__, ctx=test)

        annotations = Annotation.get_local_annotations(test.__init__, ctx=test)

        self.assertTrue(annotations)

    def test_dboundconstructor(self):
        """Test to annotate directly a bound constructor.
        """

        class Test:
            @Annotation()
            def __init__(self):
                pass

        test = Test()

        annotations = Annotation.get_local_annotations(test.__init__, ctx=test)

        self.assertTrue(annotations)


class GetParameterizedAnnotationsTest(AnnotationTest):

    def setUp(self):

        super(GetParameterizedAnnotationsTest, self).setUp()

        class BaseTest:
            pass

        class Test(BaseTest):
            pass

        self.annotation(BaseTest)
        self.annotation(Test)

        self.Test = Test
        self.BaseTest = BaseTest

        annotations = Annotation.get_annotations(BaseTest)

        self.assertEqual(len(annotations), 1)

        annotations = Annotation.get_annotations(Test)

        self.assertEqual(len(annotations), 2)

    def test_override(self):
        """
        Test to override annotation
        """

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 2)

        self.annotation.override = True

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 0)

    def test_propagate(self):
        """
        Test to propagate annotation
        """

        self.annotation.propagate = False

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

    def test_exclude(self):
        """
        Test to exclude annotations
        """

        test_annotation = TestAnnotation()

        test_annotation(self.Test)
        test_annotation(self.BaseTest)

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 4)

        annotations = Annotation.get_annotations(
            self.Test, exclude=TestAnnotation)

        self.assertEqual(len(annotations), 2)

    def test_stop_propagation(self):
        """
        Test Stop propagation annotation
        """

        stop_propagation = StopPropagation(Annotation)

        stop_propagation(self.BaseTest)

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

        stop_propagation.__del__()

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 2)

        stop_propagation(self.Test)

        annotations = Annotation.get_annotations(self.Test)

        self.assertEqual(len(annotations), 0)


class GetLocalParameterizedAnnotationsTest(AnnotationTest):

    def setUp(self):

        super(GetLocalParameterizedAnnotationsTest, self).setUp()

        class BaseTest:
            pass

        class Test(BaseTest):
            pass

        self.annotation(BaseTest)
        self.annotation(Test)

        self.Test = Test
        self.BaseTest = BaseTest

        annotations = Annotation.get_local_annotations(BaseTest)

        self.assertEqual(len(annotations), 1)

        annotations = Annotation.get_local_annotations(Test)

        self.assertEqual(len(annotations), 1)

    def test_override(self):
        """
        Test to override annotation
        """

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

        self.annotation.override = True

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

    def test_propagate(self):
        """
        Test to propagate annotation
        """

        self.annotation.propagate = False

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

    def test_exclude(self):
        """
        Test to exclude annotations
        """

        test_annotation = TestAnnotation()

        test_annotation(self.Test)
        test_annotation(self.BaseTest)

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 2)

        annotations = Annotation.get_local_annotations(
            self.Test, exclude=TestAnnotation)

        self.assertEqual(len(annotations), 1)

    def test_stop_propagation(self):
        """
        Test Stop propagation annotation
        """

        stop_propagation = StopPropagation(Annotation)

        stop_propagation(self.BaseTest)

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

        stop_propagation.__del__()

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 1)

        stop_propagation(self.Test)

        annotations = Annotation.get_local_annotations(self.Test)

        self.assertEqual(len(annotations), 2)


class GetAnnotatedFields(AnnotationTest):
    """
    Test get_annotated_fields class method
    """

    def test_class(self):

        cls = object

        members = set()

        for _, member in getmembers(cls):

            try:
                self.annotation(member, ctx=cls)
            except:
                pass
            else:
                members.add(member)

        annotated_fields = Annotation.get_annotated_fields(object)

        self.assertEqual(len(annotated_fields), len(members))

        for annotated_field in annotated_fields:

            annotations = annotated_fields[annotated_field]

            self.assertIs(annotations[0], self.annotation)


class RoutineAnnotationTest(AnnotationTest):

    def test(self):

        self.routine = 'routine'
        self.params = 'params'
        self.result = 'result'

        self.annotation = RoutineAnnotation(
            routine=self.routine, params=self.params, result=self.result
        )

        self.assertEqual(self.annotation.routine, self.routine)
        self.assertEqual(self.annotation.params, self.params)
        self.assertEqual(self.annotation.result, self.result)


if __name__ == '__main__':
    main()
