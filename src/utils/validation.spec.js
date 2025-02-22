import { test } from 'node:test';
import assert from 'node:assert';
import { validatePathsAreInScope, validateFilePath } from './validation.js';

const FN_THROW_MSG = 'Function threw an error';
const FN_NOT_THROW_MSG = 'Function did not throw an error';
const FILE_100 =
  'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.js';
const PATH_100 =
  'dir/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.js';

test('validation - validateFilesAreInScope', async (t) => {
  const scope = 'src/test-scope';

  await t.test('should pass if all modified files are in scope', () => {
    const diff = ['src/test-scope'];
    assert.doesNotThrow(() => {
      validatePathsAreInScope(diff, scope);
    }, FN_THROW_MSG);
  });

  await t.test('should throw if any modified file path is not in scope', () => {
    let diff = ['out-of-scope'];
    assert.throws(
      () => {
        validatePathsAreInScope(diff, scope);
      },
      Error,
      FN_NOT_THROW_MSG
    );

    diff = ['src/test-scope/valid', 'invalid'];
    assert.throws(
      () => {
        validatePathsAreInScope(diff, scope);
      },
      Error,
      FN_NOT_THROW_MSG
    );
  });
});

test('validation - validateFilePath', async (t) => {
  await t.test('should throw if file path is more than 100 characters', () => {
    const filePaths = ['sub/' + PATH_100, 'a' + FILE_100];

    for (let i = 0; i < filePaths.length; i++) {
      const f = filePaths[i];
      assert.throws(
        () => {
          validateFilePath(f);
        },
        Error,
        FN_NOT_THROW_MSG
      );
    }
  });

  await t.test(
    'should throw if folder or file name does not stick to convention',
    () => {
      const filePaths = [
        'notKebab/correct-file.js',
        'not kebab/correct-file.js',
        ' no-space-before.js',
        'no-space-after.js ',
        '//root/sub/dir/correct.js',
        '?correct.js',
        'notKebab.js',
        'not_kebab.js',
        'Not Kebab.js',
        'inv@lid.js',
      ];

      for (let i = 0; i < filePaths.length; i++) {
        const f = filePaths[i];
        assert.throws(
          () => {
            validateFilePath(f);
          },
          Error,
          FN_NOT_THROW_MSG
        );
      }
    }
  );

  await t.test('should throw if file extension is not allowed', () => {
    const filePaths = [
      'noextension',
      'invalid.sh',
      'invalid.zsh',
      'invalid.exe',
      'invalid.txt',
      'invalid.md',
      'invalid.py',
      'invalid.made.up',
      'invalid.png',
      'invalid.jpg',
      'invalid.jpeg',
      'invalid.wav',
      'invalid.mp3',
      'invalid.mp4',
      'invalid.mov',
    ];

    for (let i = 0; i < filePaths.length; i++) {
      const f = filePaths[i];
      assert.throws(
        () => {
          validateFilePath(f);
        },
        Error,
        FN_NOT_THROW_MSG
      );
    }
  });

  await t.test('should pass if file name and extension are valid', () => {
    const filePaths = [
      'correct.js',
      'correct-as-well.json',
      'corret-page.astro',
      'corret-page-2.astro',
      'sketch-001-99.astro',
      '8.js',
      'dir/correct.js',
      'sub/dir/correct.js',
      '/root/sub/dir/correct.js',
      '1/1/1/correct.js',
      FILE_100,
      PATH_100,
    ];

    for (let i = 0; i < filePaths.length; i++) {
      const f = filePaths[i];
      assert.doesNotThrow(() => {
        validateFilePath(f);
      }, FN_THROW_MSG);
    }
  });
});
