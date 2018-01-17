Weekly Review Procedures:

1. Process Inbox - Create many new projects as necessary (active/someday/maybe)
2. Synchronize Task Completion Status across context lists and projects
3. Reject/Defer/Update/Complete Projects and Tasks (1 project at a time)
4. Generate fresh task listings sorted by context

-----------------------------------

Weekly Review:

	gtd review

Synchronize Task Completion Status across context lists and projects

	gtd tasks sync

Generate fresh task listings sorted by context (log projects that don't have any next action):

	gtd tasks sweep

List all projects (display <h1> header, or filename if not present):

	gtd project list

Review each project in turn in a REPL session combined w/ editor sessions:

	gtd project list -review

Create many projects in a REPL session:

	gtd project create-many

Specify each relevant section from CLI (only name is required):

	gtd project create -name "Hi" -outcome "Something" -next-action "Something simple" -info "Even more stuff"

Create blank:

	gtd project create -name "Hi" -blank

Create project from CLI and skip editor session:

	gtd project create -name "Hi" -static

Renegotiate project status:

	gtd project update -id 42 -someday
	gtd project update -id 42 -maybe
	gtd project update -id 42 -reject
	gtd project update -id 42 -complete

