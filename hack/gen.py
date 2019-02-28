ALL_CAPS = ('hdmi', 'rds', 'arc', 'dts', 'usb', 'dirac', 'fm', 'am')
def const_case(astr):
    tokens = [t.lower() for t in astr.split('_')]
    capped = []
    for t in tokens:
        if '.' in t:
            t = t.replace('.', '')
        if t in ALL_CAPS:
            capped.append(t.upper())
        else:
            capped.append(t.capitalize())
    return ''.join(capped)
print '''// The code in this file is generated!
// If you want to change it, update hack/gen.py'''

with open('command_tags') as infd:
    lines = [l.strip() for l in sorted(infd)]
print '''package v1

const (
\t%sCommand CommandTag = iota''' % const_case(lines[0])
for line in lines[1:]:
    print '\t%sCommand' % const_case(line)

print ''')

var CommandTagStrings = []string{'''
for line in lines:
    print '\t"%s",' % line
print '}'
print

with open('notification_tags') as infd:
    lines = [l.strip() for l in sorted(infd)]
print '''
const (
\t%sNotification NotificationTag = iota''' % const_case(lines[0])
for line in lines[1:]:
    print '\t%sNotification' % const_case(line)
print ''')

var NotificationTagStrings = []string{'''
for line in lines:
    print '\t"%s",' % line
print '}'
print
