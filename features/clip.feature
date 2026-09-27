Feature: Clip a time range from a YouTube video
  As a video editor
  I want to download only a specific time range of a YouTube video
  So that I do not have to download and trim the full file

  Background:
    Given yt-dlp and ffmpeg are available in PATH

  Scenario: Download a valid range
    When I run "cvtr --url https://youtu.be/dQw4w9WgXcQ --start 1:05 --end 1:30"
    Then yt-dlp is called with the section "*65-90"
    And the command exits with code 0

  Scenario: Preview the command without downloading
    When I run "cvtr --url https://youtu.be/dQw4w9WgXcQ --start 10 --end 20 --dry-run"
    Then the yt-dlp command is printed
    And nothing is downloaded

  Scenario Outline: Reject invalid input before calling yt-dlp
    When I run "cvtr --url <url> --start <start> --end <end>"
    Then an error message is printed
    And the command exits with code 2
    And yt-dlp is not called

    Examples:
      | url                          | start | end  |
      | https://vimeo.com/1          | 1     | 2    |
      | https://youtu.be/dQw4w9WgXcQ | abc   | 2    |
      | https://youtu.be/dQw4w9WgXcQ | 1     | 1:60 |
      | https://youtu.be/dQw4w9WgXcQ | 20    | 10   |

  Scenario: Report a missing dependency
    Given ffmpeg is not available in PATH
    When I run "cvtr --url https://youtu.be/dQw4w9WgXcQ --start 1 --end 2"
    Then the error message mentions "ffmpeg not found in PATH"
    And the command exits with code 1
