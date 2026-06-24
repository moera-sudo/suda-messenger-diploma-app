import 'package:flutter_contacts/flutter_contacts.dart';
import 'package:injectable/injectable.dart';

/// Result of a device-contacts load: whether permission was granted and the
/// contacts (empty when denied).
typedef DeviceContactsResult = ({bool granted, List<Contact> contacts});

/// Caches device phone contacts for the app session so the Contacts tab does
/// not re-read the address book (a main-thread mapping) on every visit.
@lazySingleton
class DeviceContactsService {
  List<Contact>? _cache;

  Future<DeviceContactsResult> load({bool forceReload = false}) async {
    if (!forceReload && _cache != null) {
      return (granted: true, contacts: _cache!);
    }
    final granted = await FlutterContacts.requestPermission(readonly: true);
    if (!granted) {
      return (granted: false, contacts: const <Contact>[]);
    }
    _cache = await FlutterContacts.getContacts(withProperties: true);
    return (granted: true, contacts: _cache!);
  }
}
