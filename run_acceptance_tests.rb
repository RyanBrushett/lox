#!/usr/bin/env ruby

require 'rspec'
require 'open3'

RSpec.configure do |config|
  config.formatter = :documentation
  config.color = true
  config.tty = true
end

RSpec.describe 'The lox interpreter implemented in golang' do
  PROJECT_DIR = File.expand_path(__dir__, 'lox')
  PATH_TO_INTERPRETER = File.join(PROJECT_DIR, 'lox')
  BOOK_DIR = ENV['BOOK_DIR']

  # Set chapter value to :run in order to execute that chapter's tests
  CHAPTERS = {
    chap04_scanning:    :run,
    chap06_parsing:     :todo,
    chap07_evaluating:  :todo,
    chap08_statements:  :todo,
    chap09_control:     :todo,
    chap10_functions:   :todo,
    chap11_resolving:   :todo,
    chap12_classes:     :todo,
    chap13_inheritance: :todo,
  }

  unless BOOK_DIR
    puts "Environment variable BOOK_DIR must be set"
    exit
  end

  unless File.file?(PATH_TO_INTERPRETER)
    puts "Interpreter not found at #{PATH_TO_INTERPRETER}"
    exit
  end

  CHAPTERS.select { |_, v| v == :run }.keys.each do |chapter|
    it "passes the #{chapter} test suite" do
      output, error, status =
        Open3.capture3 \
          'dart',
          'tool/bin/test.dart',
          chapter.to_s,
          '--interpreter',
          PATH_TO_INTERPRETER,
          chdir: BOOK_DIR

      failure_message = error.empty? ? (output.slice(%r{^FAIL(.+\n)+}) || output) : error
      expect(status).to(be_success, failure_message)
    end
  end
end
